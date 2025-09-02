use ratatui::{
    layout::Constraint,
    style::{Color, Modifier, Style, Stylize},
    widgets::{Block, BorderType, Borders, Row, Table, TableState},
};
use yts_movies::{Filters, Response, Yts};

use crate::elements::Focus;

#[derive(Debug)]
pub struct MovieTable {
    pub table: TableState,
    pub response: Option<Response>,
    yts: Yts<'static>,
}

impl Default for MovieTable {
    fn default() -> Self {
        let mut table = TableState::default();
        table.select_first();
        table.select_first_column();
        let yts = Yts::default();

        Self {
            table,
            response: None,
            yts,
        }
    }
}

impl MovieTable {
    const TITLE: &'static str = " YTS MOVIES ";

    pub fn own_focus(&self) -> Focus {
        Focus::MovieTable
    }

    pub fn next_focus(&self) -> Focus {
        Focus::InputBox
    }

    pub fn footer(&self) -> String {
        match self.response {
            Some(ref r) => {
                format!(
                    " {} Movie/s - Page {}/{} ",
                    r.page.total, r.page.current, r.page.of
                )
            }
            None => String::from(" 0 Movie/s - Page 0/0 "),
        }
    }

    pub async fn search(&mut self, text: &str) -> yts_movies::Result {
        let response = self
            .yts
            .search_with_filter(text, Filters::default().build())
            .await?;

        self.response = Some(response);

        Ok(())
    }

    pub async fn next_page(&mut self, text: &str) -> yts_movies::Result {
        let response = self.response.as_ref().unwrap();
        let next_page = response.page.current + 1;
        if next_page <= response.page.of {
            self.response = Some(
                self.yts
                    .search_with_filter(text, Filters::default().page(next_page).build())
                    .await?,
            );
        }

        Ok(())
    }

    pub fn response_to_rows<'a>(&self) -> Vec<Row<'a>> {
        let mut rows: Vec<Vec<String>> = Vec::new();

        let Some(response) = self.response.as_ref() else {
            return vec![];
        };

        for movie in &response.movies {
            let genres = movie
                .genres
                .iter()
                .map(|g| g.to_string())
                .collect::<Vec<String>>()
                .join("/");

            rows.push(vec![
                movie.year.to_string(),
                movie.name.clone(),
                genres,
                movie.rating.to_string(),
            ]);
        }

        rows.iter()
            .map(|item| Row::new(item.iter().cloned()))
            .collect::<Vec<_>>()
    }

    pub fn render(&mut self, border_style: Style) -> Table<'_> {
        let header = Row::new(["Year", "Name", "Genre", "Rating"])
            .style(Style::new().dark_gray().bold())
            .bottom_margin(0);

        let widths = [
            Constraint::Percentage(10),
            Constraint::Percentage(50),
            Constraint::Percentage(30),
            Constraint::Percentage(10),
        ];

        let rows = self.response_to_rows();

        Table::new(rows, widths)
            .header(header)
            .block(
                Block::default()
                    .borders(Borders::ALL)
                    .border_type(BorderType::Thick)
                    .border_style(border_style)
                    .title(Self::TITLE)
                    .title_style(Style::new().white().bold())
                    .title_alignment(ratatui::layout::Alignment::Center)
                    .title_bottom(self.footer()),
            )
            .column_spacing(1)
            .style(Style::default().fg(Color::White))
            .row_highlight_style(
                Style::default()
                    .bg(Color::DarkGray)
                    .fg(Color::Black)
                    .add_modifier(Modifier::BOLD),
            )
            .column_highlight_style(Color::Gray)
            .cell_highlight_style(
                Style::default()
                    .bg(Color::DarkGray)
                    .fg(Color::Black)
                    .add_modifier(Modifier::BOLD),
            )
            .highlight_symbol("  ")
    }
}
