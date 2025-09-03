use ratatui::{
    layout::Constraint,
    style::{Color, Modifier, Style},
    widgets::{Block, BorderType, Borders, Row, Table, TableState},
};
use yts_movies::{Filters, Page, Response, Torrent, Yts};

use crate::elements::Focus;

#[derive(Debug)]
pub struct MovieTable {
    pub table_state: TableState,
    pub response: Response,
    yts: Yts<'static>,
}

impl Default for MovieTable {
    fn default() -> Self {
        let mut table_state = TableState::default();
        table_state.select_first();
        table_state.select_first_column();

        Self {
            table_state,
            response: Response {
                page: Page {
                    current: 0,
                    of: 0,
                    total: 0,
                },
                movies: vec![],
            },
            yts: Yts::default(),
        }
    }
}

impl MovieTable {
    const TITLE: &'static str = " YTS MOVIES ";

    pub fn footer(&self) -> String {
        let page = &self.response.page;
        if page.total != 0 {
            format!(
                " {} Movie/s - Page {}/{} ",
                page.total, page.current, page.of
            )
        } else {
            String::from(" 0 Movie/s - Page 0/0 ")
        }
    }

    pub async fn search(&mut self, text: &str) -> yts_movies::Result {
        let response = self
            .yts
            .search_with_filter(text, Filters::default().build())
            .await?;

        self.response = response;

        Ok(())
    }

    pub async fn search_torrents_by_movie(
        &mut self,
        index: usize,
    ) -> yts_movies::Result<Vec<Torrent>> {
        self.yts.torrents(&self.response.movies[index]).await
    }

    pub async fn next_page(&mut self, text: &str) -> yts_movies::Result {
        let response = &self.response;
        let next_page = response.page.current + 1;
        if next_page <= response.page.of {
            self.response = self
                .yts
                .search_with_filter(text, Filters::default().page(next_page).build())
                .await?;
        }

        Ok(())
    }

    fn response_to_rows<'a>(&self) -> Vec<Row<'a>> {
        let mut rows: Vec<Vec<String>> = Vec::new();

        if self.response.page.total == 0 {
            return vec![];
        };

        for movie in &self.response.movies {
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

    pub fn render(&mut self, focus: &Focus) -> (Table<'_>, Constraint) {
        let rows = self.response_to_rows();

        let (header, constraint) = if !rows.is_empty() {
            (
                ["Year", "Name", "Genre", "Rating"],
                Constraint::Length(rows.len() as u16 + 4),
            )
        } else {
            (["", "", "", ""], Constraint::Length(2))
        };

        let header = Row::new(header)
            .style(Style::new().dark_gray().bold())
            .bottom_margin(0);

        let widths = [
            Constraint::Percentage(10),
            Constraint::Percentage(50),
            Constraint::Percentage(30),
            Constraint::Percentage(10),
        ];

        let border_style = if matches!(focus, Focus::MovieTable) {
            Style::default()
                .fg(Color::Cyan)
                .add_modifier(Modifier::BOLD)
        } else {
            Style::default().fg(Color::DarkGray)
        };

        (
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
                .highlight_symbol("  "),
            constraint,
        )
    }
}
