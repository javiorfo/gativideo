use ratatui::{
    layout::Constraint,
    style::{Color, Modifier, Style},
    widgets::{Block, BorderType, Borders, Row, Table, TableState},
};
use transmission_rpc::{
    TransClient,
    types::{
        Nothing, RpcResponse, Torrent, TorrentAction, TorrentAddArgs, TorrentAddedOrDuplicate,
        TorrentGetField, TorrentStatus,
    },
};

use crate::elements::Focus;

pub struct Transmission {
    pub client: TransClient,
    pub table_state: TableState,
    pub torrents: Vec<Torrent>,
    download_dir: String,
}

impl Transmission {
    pub fn new(url: &'static str, dowload_dir: &str) -> Self {
        let mut table_state = TableState::default();
        table_state.select_first();
        table_state.select_first_column();

        Self {
            client: TransClient::new(url.parse().expect("Could not parse transmission url")),
            table_state,
            download_dir: dowload_dir.to_string(),
            torrents: Vec::new(),
        }
    }

    pub fn is_visible(&self) -> bool {
        !self.torrents.is_empty()
    }

    pub async fn add(&mut self, torrent_url: &str) -> transmission_rpc::types::Result<bool> {
        let add: TorrentAddArgs = TorrentAddArgs {
            filename: Some(torrent_url.replace(" ", "%20").to_string()),
            download_dir: Some(self.download_dir.clone()),
            ..TorrentAddArgs::default()
        };
        let res: RpcResponse<TorrentAddedOrDuplicate> = self.client.torrent_add(add).await?;

        self.scan().await;

        Ok(res.is_ok())
    }

    pub async fn toggle(&mut self, index: usize) -> transmission_rpc::types::Result<bool> {
        let torrent = &self.torrents[index];
        let torrent_action = if matches!(torrent.status.as_ref().unwrap(), TorrentStatus::Stopped) {
            TorrentAction::Start
        } else {
            TorrentAction::Stop
        };

        let id = torrent.id().as_ref().unwrap().clone();

        let res: RpcResponse<Nothing> =
            self.client.torrent_action(torrent_action, vec![id]).await?;

        Ok(res.is_ok())
    }

    pub async fn remove(&mut self, index: usize) -> transmission_rpc::types::Result<bool> {
        let torrent = &self.torrents[index];
        let id = torrent.id().as_ref().unwrap().clone();

        let res: RpcResponse<Nothing> = self.client.torrent_remove(vec![id], false).await?;

        Ok(res.is_ok())
    }

    pub async fn scan(&mut self) {
        let torrents = self
            .client
            .torrent_get(
                Some(vec![
                    TorrentGetField::Id,
                    TorrentGetField::Name,
                    TorrentGetField::PercentDone,
                    TorrentGetField::SizeWhenDone,
                    TorrentGetField::PeersSendingToUs,
                    TorrentGetField::PeersConnected,
                    TorrentGetField::IsStalled,
                    TorrentGetField::Status,
                ]),
                None,
            )
            .await
            .unwrap();

        self.torrents.clear();
        self.torrents = torrents.arguments.torrents;
    }

    pub fn render(&mut self, focus: &Focus) -> (Table<'_>, Constraint) {
        let widths = [
            Constraint::Percentage(30),
            Constraint::Percentage(10),
            Constraint::Percentage(20),
            Constraint::Percentage(20),
            Constraint::Percentage(20),
        ];

        let header = Row::new(["Name", "Size", "Downloaded", "Status", "Peers/Seeds"])
            .style(Style::new().dark_gray().bold())
            .bottom_margin(0);

        let mut rows: Vec<Vec<String>> = Vec::new();

        for torrent in &self.torrents {
            let status = if *torrent.is_stalled.as_ref().unwrap() {
                String::from("  Stalled")
            } else if *torrent.percent_done.as_ref().unwrap() == 1.0 {
                String::from("󰸞  Finished")
            } else {
                String::from("  Downloading")
            };

            rows.push(vec![
                torrent.name.as_ref().unwrap().clone(),
                format!(
                    "{:.2}GB",
                    torrent
                        .size_when_done
                        .as_ref()
                        .map_or(0.0, |&p| p as f64 / 1024.0 / 1024.0 / 1024.0)
                ),
                format!(
                    "{:.2}%",
                    torrent.percent_done.as_ref().map_or(0.0, |p| p * 100.0)
                ),
                status,
                format!(
                    "{}/{}",
                    torrent.peers_sending_to_us.as_ref().unwrap(),
                    torrent.peers_connected.as_ref().unwrap()
                ),
            ]);
        }

        let rows = rows
            .iter()
            .map(|item| Row::new(item.iter().cloned()))
            .collect::<Vec<_>>();

        let constraint = if rows.len() < 6 {
            Constraint::Length(rows.len() as u16 + 4)
        } else {
            Constraint::Length(10)
        };

        let border_style = if matches!(focus, Focus::TorrentTable) {
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
                        .borders(Borders::ALL)
                        .border_type(BorderType::Thick)
                        .border_style(border_style)
                        .title(" Downloads ")
                        .title_style(Style::new().white().bold())
                        .title_alignment(ratatui::layout::Alignment::Center),
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
                .highlight_symbol(" "),
            constraint,
        )
    }
}
