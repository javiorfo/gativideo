use ratatui::{
    layout::{Constraint, Flex, Layout, Rect},
    style::{Color, Modifier, Style},
    widgets::{Block, BorderType, Borders, Row, ScrollbarState, Table, TableState},
};
use yts_movies::{Movie, Torrent, Yts};

pub struct Popup<'a> {
    pub table_state: TableState,
    pub scroll_state: ScrollbarState,
    pub show: bool,
    title: &'a str,
    footer: &'a str,
}

impl<'a> Popup<'a> {
    pub fn new(title: &'a str, footer: &'a str) -> Popup<'a> {
        let mut table_state = TableState::default();
        table_state.select_first();
        table_state.select_first_column();

        Self {
            title,
            footer,
            table_state,
            scroll_state: ScrollbarState::default().position(1),
            show: false,
        }
    }

    pub fn centered_area(&self, area: Rect, x: u16, y: u16) -> Rect {
        let vertical = Layout::vertical([Constraint::Length(y)]).flex(Flex::Center);
        let horizontal = Layout::horizontal([Constraint::Length(x)]).flex(Flex::Center);
        let [area] = area.layout(&vertical);
        let [area] = area.layout(&horizontal);
        area
    }

    pub fn scroll_bar_up(&mut self) {
        let position = self.scroll_state.get_position();
        if position > 1 {
            self.scroll_state = self.scroll_state.position(position.saturating_sub(1));
        }
    }

    pub fn scroll_bar_down(&mut self, len: usize) {
        let position = self.scroll_state.get_position();
        if position < len - 1 {
            self.scroll_state = self.scroll_state.position(position.saturating_add(1));
        }
    }
}

pub struct PopupTorrent<'a> {
    pub popup: Popup<'a>,
    pub torrents: Vec<Torrent>,
    yts: Yts<'a>,
}

impl<'a> PopupTorrent<'a> {
    pub fn new(title: &'a str, footer: &'a str) -> PopupTorrent<'a> {
        Self {
            popup: Popup::new(title, footer),
            yts: Yts::default(),
            torrents: vec![],
        }
    }

    pub fn area(&self, area: Rect) -> Rect {
        self.popup.centered_area(area, 70, 5)
    }

    pub async fn search_torrents(&mut self, movie: &Movie) -> yts_movies::Result {
        self.torrents = self.yts.torrents(movie).await?;
        Ok(())
    }

    pub fn content(&self) -> Table<'a> {
        let widths = [
            Constraint::Percentage(15),
            Constraint::Percentage(15),
            Constraint::Percentage(30),
            Constraint::Percentage(20),
            Constraint::Percentage(20),
        ];

        let header = Row::new(["Quality", "Size", "Language", "Runtime", "Peers/Seeds"])
            .style(Style::new().dark_gray().bold())
            .bottom_margin(0);

        let mut rows: Vec<Vec<String>> = Vec::new();

        for torrent in &self.torrents {
            let quality = &torrent.quality;
            let quality: &str = quality.into();

            rows.push(vec![
                quality.to_owned(),
                torrent.size.clone(),
                torrent.language.clone(),
                torrent.runtime.clone(),
                torrent.peers_seeds.clone(),
            ]);
        }

        let rows = rows
            .iter()
            .map(|item| Row::new(item.iter().cloned()))
            .collect::<Vec<_>>();

        Table::new(rows, widths)
            .header(header)
            .block(
                Block::default()
                    .borders(Borders::ALL)
                    .border_type(BorderType::Plain)
                    .title(self.popup.title)
                    .title_style(Style::new().white().bold())
                    .title_alignment(ratatui::layout::Alignment::Center)
                    .title_bottom(self.popup.footer),
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
            .highlight_symbol(" ")
    }
}
