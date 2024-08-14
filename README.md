# nodebro

When you want to quickly take a glance at how several github based projects are developing and peruse through their release notes. vim navigational TUI git release tracking viewer. Simply navigate, (j/up, k/down) most recent release tags updated, and enter to view tag specific release note. ASC sort in order of age of release.

## install 

### download
` git clone https://github.com/jonaburg/nodebro.git `
### enter downloaded directory and build
` cd nodebro`
` go build -o nodebro `
### move binary to path
` mv nodebro /usr/local/bin/nodebro`


## Setup

there's an example.config.toml that should be populated before the script runs. include your PAT from github in order to not get rate limited for such simple queries :) 

```
mkdir -p ~/.config/nodebro
cp example.config.toml ~/.config/nodebro/config
```

or just do it manually and make a config like:

```
pat = "ghp_n8Jv8sTTkdsameMmlC2d8Enfnskklp3oeckK" // example PAT
[[repos]]
Owner = "bitcoin"
Repo = "bitcoin"
#[[repos]]
Owner = "ethereum"
Repo = "go-ethereum"
```

etc.

## run
`nodebro`

