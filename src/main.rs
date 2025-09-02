use color_eyre::Result;
use crossterm::event::{self, KeyCode};
use ratatui::Frame;
use ratatui::layout::{Constraint, Layout};
use ratatui::widgets::Clear;

mod elements;

use elements::Focus;

use crate::elements::{InputBox, MovieTable, Popup};

#[tokio::main]
async fn main() -> Result<()> {
    color_eyre::install()?;

    let mut terminal = ratatui::init();

    let mut show_popup = false;
    let mut focus = Focus::default();
    let mut input_box = InputBox::default();
    let mut movie_table = MovieTable::default();

    loop {
        terminal.draw(|frame| render(frame, &mut movie_table, &focus, &input_box, show_popup))?;
        if let Some(key) = event::read()?.as_key_press_event() {
            match focus {
                Focus::InputBox => match key.code {
                    KeyCode::Tab => {
                        focus = Focus::MovieTable;
                    }
                    KeyCode::Enter => {
                        movie_table.search(&input_box.text).await.expect("");
                        focus = Focus::MovieTable;
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
                Focus::MovieTable => match key.code {
                    KeyCode::Tab => {
                        focus = Focus::InputBox;
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
                        let _selected = movie_table.table.selected().unwrap();
                        //                         println!("{}", response.movies[selected].name);
                    }
                    KeyCode::Char('t') => {
                        show_popup = true;
                        focus = Focus::TorrentPopup;
                    }
                    _ => {}
                },

                Focus::TorrentTable => match key.code {
                    _ => {}
                },
                Focus::TorrentPopup => match key.code {
                    KeyCode::Char('q') | KeyCode::Esc => {
                        show_popup = false;
                        focus = Focus::MovieTable;
                    }
                    _ => {}
                },
            }
        }
    }
}

fn render(
    frame: &mut Frame,
    movie_table: &mut MovieTable,
    focus: &Focus,
    input_box: &InputBox,
    show_popup_torrents: bool,
) {
    let mut table_state = movie_table.table;
    let (table, constraint) = movie_table.render(focus);

    let area = frame.area();
    let layout = Layout::vertical([Constraint::Length(3), constraint]);
    let [input_box_area, movie_table_area] = area.layout(&layout);

    frame.render_widget(input_box.render(focus), input_box_area);

    if matches!(focus, Focus::InputBox) {
        frame.set_cursor_position((
            input_box_area.x + input_box.text.len() as u16 + 1,
            input_box_area.y + 1,
        ));
    }

    frame.render_stateful_widget(table, movie_table_area, &mut table_state);

    if show_popup_torrents {
        let popup = Popup::new(" Torrents ");
        let popup_area = popup.centered_area(area, 60, 40);
        frame.render_widget(Clear, popup_area);
        frame.render_widget(popup.table(), popup_area);
    }
}
