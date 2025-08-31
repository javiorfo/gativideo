use color_eyre::Result;
use crossterm::event::{self, KeyCode};
use ratatui::Frame;
use ratatui::layout::{Constraint, Direction, Layout, Rect};
use ratatui::style::{Color, Modifier, Style, Stylize};
use ratatui::widgets::{Block, BorderType, Borders, Paragraph, Row, Table, TableState};
use yts_movies::Filters;

#[tokio::main]
async fn main() -> Result<()> {
    color_eyre::install()?;

    let mut table_state = TableState::default();
    table_state.select_first();
    table_state.select_first_column();
    let mut terminal = ratatui::init();
    let mut focus = Focus::Table;
    let mut input_text = String::new();

    let yts = yts_movies::Yts::default();

    let response: yts_movies::Response = yts
        .search_with_filter("", Filters::default().build())
        .await
        .expect("error");

    let mut rows: Vec<Vec<String>> = Vec::new();
    for movie in response.movies {
        let genres = movie
            .genres
            .iter()
            .map(|g| g.to_string())
            .collect::<Vec<String>>()
            .join("/");
        rows.push(vec![
            movie.year.to_string(),
            movie.name,
            genres,
            movie.rating.to_string(),
        ]);
    }

    //     let mut rows: Vec<Vec<&str>> = vec![
    //         vec!["1972", "The Godfather", "Drama/Crime", "9.3"],
    //         vec!["2013", "Django Unchained", "Western/Drama/Crime", "8.3"],
    //         vec!["1999", "Fight Club", "Action/Drama/Comedy", "8.7"],
    //         vec!["1994", "Dumb And Dumber", "Comedy", "7.7"],
    //         vec!["1994", "Dumb And Dumber", "Comedy", "7.7"],
    //     ];

    loop {
        terminal.draw(|frame| render(frame, &mut table_state, &focus, &input_text, &rows))?;
        if let Some(key) = event::read()?.as_key_press_event() {
            match focus {
                Focus::Table => match key.code {
                    KeyCode::Tab => {
                        focus = Focus::Input;
                    }
                    KeyCode::Char('q') | KeyCode::Esc => {
                        ratatui::restore();
                        return Ok(());
                    }
                    KeyCode::Char('j') | KeyCode::Down => table_state.select_next(),
                    KeyCode::Char('k') | KeyCode::Up => table_state.select_previous(),
                    KeyCode::Char('l') | KeyCode::Right => table_state.select_next_column(),
                    KeyCode::Char('h') | KeyCode::Left => table_state.select_previous_column(),
                    KeyCode::Char('g') => table_state.select_first(),
                    KeyCode::Char('G') => table_state.select_last(),
                    _ => {}
                },
                Focus::Input => match key.code {
                    KeyCode::Tab => {
                        focus = Focus::Table;
                    }
                    KeyCode::Enter => {
                        let response: yts_movies::Response = yts
                            .search_with_filter(&input_text, Filters::default().build())
                            .await
                            .expect("error");

                        rows.clear();
                        
                        for movie in response.movies {
                            let genres = movie
                                .genres
                                .iter()
                                .map(|g| g.to_string())
                                .collect::<Vec<String>>()
                                .join("/");
                            rows.push(vec![
                                movie.year.to_string(),
                                movie.name,
                                genres,
                                movie.rating.to_string(),
                            ]);
                        }

                        focus = Focus::Table;
                    }
                    KeyCode::Char(c) => {
                        input_text.push(c);
                    }
                    KeyCode::Backspace => {
                        input_text.pop();
                    }
                    KeyCode::Esc => {
                        ratatui::restore();
                        return Ok(());
                    }
                    _ => {}
                },
            }
        }
    }
}

// Add this before the main function
enum Focus {
    Table,
    Input,
}

/// Render the UI with a table.
fn render(
    frame: &mut Frame,
    table_state: &mut TableState,
    focus: &Focus,
    input_text: &str,
    rows: &[Vec<String>],
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

    let input_border_style = if matches!(focus, Focus::Input) {
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

    let input_box = Paragraph::new(input_text)
        .style(Style::default().fg(Color::Yellow))
        .block(input_block);
    frame.render_widget(input_box, top);

    if matches!(focus, Focus::Input) {
        frame.set_cursor_position((chunks[0].x + input_text.len() as u16 + 1, chunks[0].y + 1));
    }

    // Determine the border style for the table based on focus
    let table_border_style = if matches!(focus, Focus::Table) {
        Style::default()
            .fg(Color::Cyan)
            .add_modifier(Modifier::BOLD)
    } else {
        Style::default().fg(Color::DarkGray)
    };
    //     frame.render_widget(title_paragraph, top);

    render_table(frame, main, table_state, table_border_style, rows);
}

/// Render a table with some rows and columns.
pub fn render_table(
    frame: &mut Frame,
    area: Rect,
    table_state: &mut TableState,
    border_style: Style,
    rows: &[Vec<String>],
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

    let rows: Vec<Row> = rows
        .iter()
        .map(|item| Row::new(item.iter().cloned()))
        .collect();

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
                .title_bottom(" 0 Movie/s - Page 1/0 "),
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
    //         .footer(
    //             Row::new(["0 Movie/s - Page 1/0"]).style(
    //                 Style::default()
    //                     .fg(Color::DarkGray)
    //                     .add_modifier(Modifier::BOLD),
    //             ),
    //         );

    frame.render_stateful_widget(table, area, table_state);
}
