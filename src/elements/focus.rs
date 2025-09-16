#[derive(Debug, PartialEq, Default)]
pub enum Focus {
    #[default]
    InputBox,
    MovieTable,
    TorrentTable,
    PopupTorrent,
    PopupSubtitle,
}
