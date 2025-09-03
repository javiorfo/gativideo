use color_eyre::Result;
use crossterm::event::{self, KeyCode};
use ratatui::Frame;
use ratatui::layout::{Constraint, Layout};
use ratatui::widgets::Clear;

mod elements;

use elements::Focus;

use crate::elements::{InputBox, MovieTable, PopupTorrent};

#[tokio::main]
async fn main() -> Result<()> {
    color_eyre::install()?;

    let mut terminal = ratatui::init();

    let mut focus = Focus::default();
    let mut input_box = InputBox::default();
    let mut movie_table = MovieTable::default();
    let mut popup_torrent = PopupTorrent::new(" Torrents ", " Enter to download ");

    loop {
        terminal
            .draw(|frame| render(frame, &mut movie_table, &focus, &input_box, &popup_torrent))?;
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
                    KeyCode::Char('j') | KeyCode::Down => movie_table.table_state.select_next(),
                    KeyCode::Char('k') | KeyCode::Up => movie_table.table_state.select_previous(),
                    KeyCode::Char('l') | KeyCode::Right => {
                        movie_table.next_page(&input_box.text).await.expect("");
                    }
                    KeyCode::Char('h') | KeyCode::Left => {
                        movie_table.table_state.select_previous_column()
                    }
                    KeyCode::Char('g') => movie_table.table_state.select_first(),
                    KeyCode::Char('G') => movie_table.table_state.select_last(),
                    KeyCode::Char('t') => {
                        if let Some(selected) = movie_table.table_state.selected()
                            && !movie_table.response.movies.is_empty()
                        {
                            let movie = &movie_table.response.movies[selected];
                            popup_torrent.search_torrents(movie).await.unwrap();

                            popup_torrent.popup.show = true;
                            focus = Focus::TorrentPopup;
                        }
                    }
                    _ => {}
                },

                Focus::TorrentTable => match key.code {
                    _ => {}
                },
                Focus::TorrentPopup => match key.code {
                    KeyCode::Char('j') | KeyCode::Down => {
                        popup_torrent.popup.table_state.select_next()
                    }
                    KeyCode::Char('k') | KeyCode::Up => {
                        popup_torrent.popup.table_state.select_previous()
                    }
                    KeyCode::Char('q') | KeyCode::Esc => {
                        popup_torrent.popup.show = false;
                        focus = Focus::MovieTable;
                    }
                    KeyCode::Enter => {
                        todo!("download movie")
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
    popup_torrent: &PopupTorrent,
) {
    let mut table_state = movie_table.table_state;
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

    if popup_torrent.popup.show {
        let popup_area = popup_torrent.area(movie_table_area);
        let mut table_state = popup_torrent.popup.table_state;
        frame.render_widget(Clear, popup_area);
        frame.render_stateful_widget(popup_torrent.content(), popup_area, &mut table_state);
    }
}
