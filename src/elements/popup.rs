use ratatui::{
    layout::{Constraint, Flex, Layout, Rect},
    style::{Color, Modifier, Style},
    widgets::{Block, BorderType, Borders, Row, Table},
};

pub struct Popup<'a> {
    pub title: &'a str,
}

impl<'a> Popup<'a> {
    pub fn new(title: &'a str) -> Popup<'a> {
        Self { title }
    }

    pub fn block(&self) -> Block<'a> {
        Block::bordered()
            .title(self.title)
            .title_bottom(" Press Enter to start download ")
            .title_style(Style::new().white().bold())
            .title_alignment(ratatui::layout::Alignment::Center)
    }

    pub fn table(&self) -> Table<'a> {
        let widths = [
            Constraint::Percentage(10),
            Constraint::Percentage(50),
            Constraint::Percentage(30),
            Constraint::Percentage(10),
        ];
        let header = Row::new(["Year", "Name", "Genre", "Rating"])
            .style(Style::new().dark_gray().bold())
            .bottom_margin(0);

        Table::new([Row::new(["Year", "Name", "Genre", "Rating"])], widths)
            .header(header)
            .block(
                Block::default()
                    .borders(Borders::ALL)
                    .border_type(BorderType::Thick)
                    .title(" Torrents ")
                    .title_style(Style::new().white().bold())
                    .title_alignment(ratatui::layout::Alignment::Center)
                    .title_bottom(" footer "),
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

    pub fn centered_area(&self, area: Rect, percent_x: u16, percent_y: u16) -> Rect {
        let vertical = Layout::vertical([Constraint::Percentage(percent_y)]).flex(Flex::Center);
        let horizontal = Layout::horizontal([Constraint::Percentage(percent_x)]).flex(Flex::Center);
        let [area] = area.layout(&vertical);
        let [area] = area.layout(&horizontal);
        area
    }
}
