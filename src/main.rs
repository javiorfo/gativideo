use color_eyre::Result;
use crossterm::event::{self, KeyCode};
use ratatui::Frame;
use ratatui::layout::{Constraint, Direction, Layout};
use ratatui::style::{Color, Modifier, Style};
use ratatui::widgets::{Block, BorderType, Borders, Paragraph};

mod elements;

use elements::Focus;

use crate::elements::{InputBox, MovieTable};

#[tokio::main]
async fn main() -> Result<()> {
    color_eyre::install()?;

    let mut terminal = ratatui::init();

    let mut focus = Focus::default();
    let mut input_box = InputBox::default();
    let mut movie_table = MovieTable::default();

    loop {
        terminal.draw(|frame| render(frame, &mut movie_table, &focus, &input_box))?;
        if let Some(key) = event::read()?.as_key_press_event() {
            match focus {
                Focus::MovieTable => match key.code {
                    KeyCode::Tab => {
                        focus = movie_table.next_focus();
                    }
                    KeyCode::Char('q') | KeyCode::Esc => {
                        ratatui::restore();
                        return Ok(());
                    }
                    KeyCode::Char('j') | KeyCode::Down => movie_table.table.select_next(),
                    KeyCode::Char('k') | KeyCode::Up => movie_table.table.select_previous(),
                    KeyCode::Char('l') | KeyCode::Right => {
                        movie_table.next_page(&input_box.text).await.expect("");
                    }
                    KeyCode::Char('h') | KeyCode::Left => {
                        movie_table.table.select_previous_column()
                    }
                    KeyCode::Char('g') => movie_table.table.select_first(),
                    KeyCode::Char('G') => movie_table.table.select_last(),
                    KeyCode::Enter => {
                        let selected = movie_table.table.selected().unwrap();
                        //                         println!("{}", response.movies[selected].name);
                    }
                    _ => {}
                },
                Focus::InputBox => match key.code {
                    KeyCode::Tab => {
                        focus = input_box.next_focus();
                    }
                    KeyCode::Enter => {
                        movie_table.search(&input_box.text).await.expect("");
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

fn render(frame: &mut Frame, movie_table: &mut MovieTable, focus: &Focus, input_box: &InputBox) {
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

    let table_border_style = if focus == &movie_table.own_focus() {
        Style::default()
            .fg(Color::Cyan)
            .add_modifier(Modifier::BOLD)
    } else {
        Style::default().fg(Color::DarkGray)
    };

    let mut table_state = movie_table.table.clone();
    let table = movie_table.render(table_border_style);
    frame.render_stateful_widget(table, main, &mut table_state);
}
