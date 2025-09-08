use ratatui::{
    layout::Constraint,
    style::{Color, Modifier, Style},
    widgets::{Block, BorderType, Borders, Row, Table, TableState},
};
use transmission_rpc::{
    TransClient,
    types::{RpcResponse, Torrent, TorrentAddArgs, TorrentAddedOrDuplicate, TorrentGetField},
};

use crate::elements::Focus;

pub struct Transmission {
    pub client: TransClient,
    pub table_state: TableState,
    pub torrents: Vec<Torrent>,
    download_dir: String,
}

impl Transmission {
    pub fn new(url: &str, dowload_dir: &str) -> Self {
        let mut table_state = TableState::default();
        table_state.select_first();
        table_state.select_first_column();

        Self {
            client: TransClient::new(url.parse().expect("Could not parse transmission url")),
            table_state,
            download_dir: dowload_dir.to_string(),
            torrents: vec![],
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

        self.scan().await.unwrap();

        Ok(res.is_ok())
    }

    pub async fn scan(&mut self) -> transmission_rpc::types::Result<()> {
        let torrents = self
            .client
            .torrent_get(
                Some(vec![
                    TorrentGetField::Name,
                    TorrentGetField::PercentDone,
                    TorrentGetField::SizeWhenDone,
                    TorrentGetField::PeersSendingToUs,
                    TorrentGetField::PeersConnected,
                    TorrentGetField::IsFinished,
                    TorrentGetField::IsStalled,
                ]),
                None,
            )
            .await?;

        //         tokio::time::sleep(tokio::time::Duration::from_secs(3)).await;

        self.torrents = torrents.arguments.torrents;
        Ok(())
    }

    pub fn render(&mut self, focus: &Focus) -> Table<'_> {
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
            } else if *torrent.is_finished.as_ref().unwrap() {
                String::from("󰸞  Finished")
            } else {
                String::from("  Downloading")
            };

            rows.push(vec![
                torrent.name.as_ref().unwrap().clone(),
                format!(
                    "{:.2}%",
                    torrent
                        .size_when_done
                        .as_ref()
                        .map_or(0.0, |&p| p as f64 / 1024.0 / 1024.0 / 1024.0)
                ),
                format!(
                    "{:.2}GB",
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

        let border_style = if matches!(focus, Focus::TorrentTable) {
            Style::default()
                .fg(Color::Cyan)
                .add_modifier(Modifier::BOLD)
        } else {
            Style::default().fg(Color::DarkGray)
        };

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
            .highlight_symbol(" ")
    }
}
