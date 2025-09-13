use color_eyre::Result;
use crossterm::event::{self, KeyCode};
use ratatui::Frame;
use ratatui::layout::{Constraint, Layout};
use ratatui::symbols::scrollbar;
use ratatui::widgets::{Clear, Scrollbar, ScrollbarOrientation};

use crate::config::configuration;
use crate::downloads::Transmission;
use crate::elements::{Focus, InputBox, MovieTable, PopupTorrent};

pub async fn run() -> Result<()> {
    let config = configuration();

    color_eyre::install()?;

    let mut terminal = ratatui::init();

    let mut focus = Focus::default();
    let mut input_box = InputBox::default();
    let mut movie_table = MovieTable::default();
    let mut popup_torrent = PopupTorrent::new(" Torrents ", " Enter to download ");
    let mut transmission = Transmission::new(config.transmission_host, config.yts_download_dir);

    let mut last_redraw_time = tokio::time::Instant::now();
    let redraw_interval = tokio::time::Duration::from_secs(1);

    transmission.scan().await;

    loop {
        terminal.draw(|frame| {
            render(
                frame,
                &mut movie_table,
                &focus,
                &input_box,
                &popup_torrent,
                &mut transmission,
            )
        })?;

        let time_since_last_redraw = tokio::time::Instant::now().duration_since(last_redraw_time);
        let timeout = redraw_interval.saturating_sub(time_since_last_redraw);
        if tokio::time::Instant::now().duration_since(last_redraw_time) >= redraw_interval {
            transmission.scan().await;
            last_redraw_time = tokio::time::Instant::now();
        }

        if event::poll(timeout)?
            && let Some(key) = event::read()?.as_key_press_event()
        {
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
                        focus = if transmission.is_visible() {
                            Focus::TorrentTable
                        } else {
                            Focus::InputBox
                        };
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
                        movie_table.previous_page(&input_box.text).await.expect("");
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
                    KeyCode::Char('s') => {
                        if let Some(selected) = transmission.table_state.selected()
                            && !transmission.torrents.is_empty()
                        {
                            transmission.toggle(selected).await.unwrap();
                        }
                    }
                    KeyCode::Char('q') | KeyCode::Esc => {
                        ratatui::restore();
                        return Ok(());
                    }
                    KeyCode::Tab => {
                        focus = Focus::InputBox;
                    }
                    KeyCode::Char('j') | KeyCode::Down => {
                        transmission.table_state.select_next();
                        transmission.scroll_bar_up();
                    }
                    KeyCode::Char('k') | KeyCode::Up => {
                        transmission.table_state.select_previous();
                        transmission.scroll_bar_down();
                    }
                    KeyCode::Char('d') => {
                        if let Some(selected) = transmission.table_state.selected()
                            && !transmission.torrents.is_empty()
                        {
                            transmission.remove(selected).await.unwrap();
                        }
                    }
                    _ => {}
                },
                Focus::TorrentPopup => match key.code {
                    KeyCode::Char('j') | KeyCode::Down => {
                        popup_torrent.popup.table_state.select_next();
                        popup_torrent
                            .popup
                            .scroll_bar_down(popup_torrent.torrents.len());
                    }
                    KeyCode::Char('k') | KeyCode::Up => {
                        popup_torrent.popup.table_state.select_previous();
                        popup_torrent.popup.scroll_bar_up();
                    }
                    KeyCode::Char('q') | KeyCode::Esc => {
                        popup_torrent.popup.show = false;
                        focus = Focus::MovieTable;
                    }
                    KeyCode::Enter => {
                        if let Some(selected) = popup_torrent.popup.table_state.selected() {
                            let torrent = &popup_torrent.torrents[selected];
                            transmission.add(&torrent.link).await.unwrap();
                        }
                        popup_torrent.popup.show = false;
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
    popup_torrent: &PopupTorrent,
    transmission: &mut Transmission,
) {
    let mut movie_table_state = movie_table.table_state;
    let (table, constraint) = movie_table.render(focus);

    let visible = transmission.is_visible();
    let mut transmission_table_state = transmission.table_state;
    let (torrent_table, torrent_constraint) = transmission.render(focus);

    let area = frame.area();
    let layout = Layout::vertical([Constraint::Length(3), constraint, torrent_constraint]);
    let [input_box_area, movie_table_area, torrent_table_area] = area.layout(&layout);

    frame.render_widget(input_box.render(focus), input_box_area);

    if matches!(focus, Focus::InputBox) {
        frame.set_cursor_position((
            input_box_area.x + input_box.text.len() as u16 + 1,
            input_box_area.y + 1,
        ));
    }

    frame.render_stateful_widget(table, movie_table_area, &mut movie_table_state);

    if popup_torrent.popup.show {
        let popup_area = popup_torrent.area(movie_table_area);
        let mut table_state = popup_torrent.popup.table_state;
        frame.render_widget(Clear, popup_area);
        frame.render_stateful_widget(popup_torrent.content(), popup_area, &mut table_state);

        let mut scroll_state = popup_torrent
            .popup
            .scroll_state
            .content_length(popup_torrent.torrents.len() + 2);

        frame.render_stateful_widget(
            Scrollbar::new(ScrollbarOrientation::VerticalRight)
                .symbols(scrollbar::VERTICAL)
                .begin_symbol(None)
                .track_symbol(None)
                .end_symbol(None),
            popup_area,
            &mut scroll_state,
        );
    }

    if visible {
        frame.render_stateful_widget(
            torrent_table,
            torrent_table_area,
            &mut transmission_table_state,
        );

        let torrents_len = transmission.torrents.len();

        if torrents_len > 5 {
            let mut scroll_state = transmission.scroll_state.content_length(torrents_len + 2);

            frame.render_stateful_widget(
                Scrollbar::new(ScrollbarOrientation::VerticalRight)
                    .symbols(scrollbar::VERTICAL)
                    .begin_symbol(None)
                    .track_symbol(None)
                    .end_symbol(None),
                torrent_table_area,
                &mut scroll_state,
            );
        }
    }
}
