use color_eyre::Result;
use crossterm::event::{self, KeyCode};
use ratatui::Frame;
use ratatui::layout::{Constraint, Direction, Layout, Rect};
use ratatui::style::{Color, Modifier, Style, Stylize};
use ratatui::widgets::{Block, BorderType, Borders, Paragraph, Row, Table, TableState};
use yts_movies::Filters;

mod elements;

use elements::Focus;

use crate::elements::InputBox;

#[tokio::main]
async fn main() -> Result<()> {
    color_eyre::install()?;

    let mut table_state = TableState::default();
    table_state.select_first();
    table_state.select_first_column();
    let mut terminal = ratatui::init();

    let mut focus = Focus::default();
    let mut input_box = InputBox::default();

    let yts = yts_movies::Yts::default();

    let mut response: yts_movies::Response = yts
        .search_with_filter("", Filters::default().build())
        .await
        .expect("error");

    loop {
        terminal.draw(|frame| render(frame, &mut table_state, &focus, &input_box, &response))?;
        if let Some(key) = event::read()?.as_key_press_event() {
            match focus {
                Focus::MovieTable => match key.code {
                    KeyCode::Tab => {
                        focus = Focus::InputBox;
                    }
                    KeyCode::Char('q') | KeyCode::Esc => {
                        ratatui::restore();
                        return Ok(());
                    }
                    KeyCode::Char('j') | KeyCode::Down => table_state.select_next(),
                    KeyCode::Char('k') | KeyCode::Up => table_state.select_previous(),
                    KeyCode::Char('l') | KeyCode::Right => {
                        let next_page = response.page.current + 1;
                        if next_page < response.page.of {
                            response = yts
                                .search_with_filter(
                                    &input_box.text,
                                    Filters::default().page(next_page).build(),
                                )
                                .await
                                .expect("error");
                        }
                    }
                    KeyCode::Char('h') | KeyCode::Left => table_state.select_previous_column(),
                    KeyCode::Char('g') => table_state.select_first(),
                    KeyCode::Char('G') => table_state.select_last(),
                    KeyCode::Enter => {
                        let selected = table_state.selected().unwrap();
                        println!("{}", response.movies[selected].name);
                    }
                    _ => {}
                },
                Focus::InputBox => match key.code {
                    KeyCode::Tab => {
                        focus = input_box.next_focus();
                    }
                    KeyCode::Enter => {
                        response = yts
                            .search_with_filter(&input_box.text, Filters::default().build())
                            .await
                            .expect("error");

                        focus = input_box.next_focus();
                    }
                    KeyCode::Char(c) => {
                        input_box.text.push(c);
                    }
                    KeyCode::Backspace => {
                        input_box.text.pop();
                    }
                    KeyCode::Esc => {
                        ratatui::restore();
                        return Ok(());
                    }
                    _ => {}
                },
                Focus::TorrentTable => match key.code {
                    _ => {}
                },
            }
        }
    }
}

fn render(
    frame: &mut Frame,
    table_state: &mut TableState,
    focus: &Focus,
    input_box: &InputBox,
    response: &yts_movies::Response,
) {
    let chunks = Layout::default()
        .direction(Direction::Vertical)
        .constraints(
            [
                Constraint::Length(3),
                Constraint::Percentage(42),
                Constraint::Min(0),
            ]
            .as_ref(),
        )
        .spacing(1)
        .split(frame.area());

    let top = chunks[0];
    let main = chunks[1];

    let input_border_style = if focus == &input_box.own_focus() {
        Style::default()
            .fg(Color::Yellow)
            .add_modifier(Modifier::BOLD)
    } else {
        Style::default().fg(Color::DarkGray)
    };

    let input_block = Block::default()
        .borders(Borders::ALL)
        .border_type(BorderType::Thick)
        .border_style(input_border_style)
        .title("Input");

    let paragraph = Paragraph::new(input_box.text.clone())
        .style(Style::default().fg(Color::Yellow))
        .block(input_block);
    frame.render_widget(paragraph, top);

    if focus == &input_box.own_focus() {
        frame.set_cursor_position((
            chunks[0].x + input_box.text.len() as u16 + 1,
            chunks[0].y + 1,
        ));
    }

    let table_border_style = if matches!(focus, Focus::MovieTable) {
        Style::default()
            .fg(Color::Cyan)
            .add_modifier(Modifier::BOLD)
    } else {
        Style::default().fg(Color::DarkGray)
    };

    render_table(frame, main, table_state, table_border_style, response);
}

pub fn render_table(
    frame: &mut Frame,
    area: Rect,
    table_state: &mut TableState,
    border_style: Style,
    response: &yts_movies::Response,
) {
    let header = Row::new(["Year", "Name", "Genre", "Rating"])
        .style(Style::new().dark_gray().bold())
        .bottom_margin(0);

    let widths = [
        Constraint::Percentage(10),
        Constraint::Percentage(50),
        Constraint::Percentage(30),
        Constraint::Percentage(10),
    ];

    let mut rows: Vec<Vec<String>> = Vec::new();
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

    let rows: Vec<Row> = rows
        .iter()
        .map(|item| Row::new(item.iter().cloned()))
        .collect();

    let footer = format!(
        " {} Movie/s - Page {}/{} ",
        response.page.total, response.page.current, response.page.of
    );

    let table = Table::new(rows, widths)
        .header(header)
        .block(
            Block::default()
                .borders(Borders::ALL)
                .border_type(BorderType::Thick)
                //                 .border_style(Style::new().dark_gray())
                .border_style(border_style)
                .title(" YTS MOVIES ")
                .title_style(Style::new().white().bold())
                .title_alignment(ratatui::layout::Alignment::Center)
                .title_bottom(footer),
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
        .highlight_symbol("  ");

    frame.render_stateful_widget(table, area, table_state);
}
