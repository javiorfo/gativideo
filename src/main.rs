use color_eyre::Result;

mod app;
mod config;
mod downloads;
mod elements;

#[tokio::main]
async fn main() -> Result<()> {
    app::run().await
}
