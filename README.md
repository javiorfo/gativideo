# gativideo
*TUI app to download YTS movies and opensubtitles subtitles*

## Caveats
- Rust version **1.88**
- It's upon `transmission-rpc` protocol. So It requires to be active in order to use gativideo. 
- This program has been developed on and for Linux following open source philosophy.

## Installation
- Using Cargo
```bash
cargo install gativideo
```

- From AUR Arch Linux:
```bash
paru -S gativideo
```

## Description
- This program is a TUI wrapper of `YTS movies (a.k.a. yify)` and `opensubtitles.com (a.k.a. opensubtitles.org)`
- It uses `transmission-rpc` protocol
- You can search and download movies from `YTS` and search and download subtitles from `opensubtitles`. 
- This program serves itself from crates [yts-movies](https://github.com/javiorfo/yts-movies) and [opensubs](https://github.com/javiorfo/opensubs)
- Some properties could be define in a file stored as `.config/gativideo/config.toml` [default values](https://github.com/javiorfo/gativideo/blob/master/example/config.toml)

## Usage
- Write the name of a movie and press `enter` to search
- Use the `up` and `down` keys to navigate the movies table or the subtitles table
- Use `l` to go to the next page
- Use `h` to go to the previous page
- Use `Ctrl+s` to toggle between movies table and subtitles table
    - The subtitle is bound to the movie (if there is no subtitle for the movie, the subtitle table will be empty)
- Use `Ctrl+d` to add a movie or a subtitle to download
    - The subtitle will be download almost instantly 
    - The movie will be added to an internal torrent client which will show the progress and peers connected
- Use `Ctrl+r` to cancel download
- Use `Ctrl+c` to exit **gativideo**

## Extra
- The subtitle search is disabled by default. In case is needed, enable it and set the language using the **config.toml** [parameter](https://github.com/javiorfo/gativideo/blob/master/example/config.toml).
- The subtitle search is bound to the movie. Nonetheless, download only the subtitles can be done.
- Multiple movies at the time can be downloaded. You can close **gativideo** and the downloads still continue.
    
## Demos and screenshots

<img src="https://github.com/javiorfo/img/blob/master/gativideo/gativideo-simple.gif?raw=true" alt="gativideo"/>

#### Using filters
- **order** filter could be: *latest, oldest, rating, alphabetical, featured, year or likes*
<img src="https://github.com/javiorfo/img/blob/master/gativideo/gativideo-filters.gif?raw=true" alt="gativideo"/>

---

### Donate
- **Bitcoin** [(QR)](https://raw.githubusercontent.com/javiorfo/img/master/crypto/bitcoin.png)  `1GqdJ63RDPE4eJKujHi166FAyigvHu5R7v`
- [Paypal](https://www.paypal.com/donate/?hosted_button_id=FA7SGLSCT2H8G)

