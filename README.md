# bitsmuggler
*TUI app to download YTS movies and opensubtitles subs*

## Caveats
- Go version **1.24**
- This program has been developed on and for Linux following open source philosophy.

<img src="https://github.com/javiorfo/img/blob/master/bitsmuggler/bitsmuggler.png?raw=true" alt="bitsmuggler"/>

## Installation
- Using Go
```bash
go install github.com/javiorfo/bitsmuggler@latest
```

- Downloading, compiling and installing manually (Linux):
```bash
git clone https://github.com/javiorfo/bitsmuggler
cd bitsmuggler
sudo make clean install
```

- From AUR Arch Linux:
```bash
yay -S bitsmuggler
```

## Description
- This program is a kind of TUI wrapper of `YTS movies (a.k.a. yify)` server and `opensubtitles.com (a.k.a. opensubtitles.org)`
- You can search and download movies from `YTS` server and search and download subtitles from `opensubtitles` server
- Some properties could be define in a file stored as `.config/bitsmuggler/config.toml` [default values](https://github.com/javiorfo/bitsmuggler/blob/master/example/config.toml)

## Usage
- When **bitsmuggler** is executed it will search the latest movies ordered by rating 
- Write the name of a movie and press `enter` to search
- Use the `up` and `down` keys to navigate the movies table or the subtitles table
- Use `Ctrl+n` to go to the next page
- Use `Ctrl+p` to go to the previous page
- Use `Ctrl+s` to toggle between movies table and subtitles table
    - The subtitle is bound to the movie (if there is no subtitle for the movie, the subtitle table will be empty)
- Use `Ctrl+d` to add a movie or a subtitle to download
    - The subtitle will be download almost instantly 
    - The movie will be added to an internal torrent client which will show the progress and peers connected
- Use `Ctrl+r` to cancel download
- Use `Ctrl+c` to exit **bitsmuggler**

## Extra
- The subtitle search is disabled by default. In case is needed, enable it and set the language using the **config.toml** [parameter](https://github.com/javiorfo/bitsmuggler/blob/master/example/config.toml).
- The subtitle search is bound to the movie. Nonetheless, you can download only the subtitles.
- Only one movie at the time could be downloaded. You can close **bitsmuggler** and when you open it again It will resume the incomplete download and finish it unless it is canceled by the user.
- The quality selected by default is **1080** but could be modified in the config.toml file. Using 2160, 1080 or 720 (A fallback is used from mayor to minor. Ex: If 2160 is set and not found, it will search 1080 and so on).
    
## Demos and screenshots

<img src="https://github.com/javiorfo/img/blob/master/bitsmuggler/bitsmuggler-simple.gif?raw=true" alt="bitsmuggler"/>

#### Using filters
- **order** filter could be: *latest, oldest, rating, alphabetical, featured, year or likes*
<img src="https://github.com/javiorfo/img/blob/master/bitsmuggler/bitsmuggler-filters.gif?raw=true" alt="bitsmuggler"/>

---

### Donate
- **Bitcoin** [(QR)](https://raw.githubusercontent.com/javiorfo/img/master/crypto/bitcoin.png)  `1GqdJ63RDPE4eJKujHi166FAyigvHu5R7v`
- [Paypal](https://www.paypal.com/donate/?hosted_button_id=FA7SGLSCT2H8G)

