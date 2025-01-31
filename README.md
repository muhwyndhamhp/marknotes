# What is this?
A simple static blog built with Golang + HTMX + Tailwind CSS. With some fancy secret under the hood. See the Video below for all the available features:

https://github.com/user-attachments/assets/cfb811f7-3329-4dfa-bc8e-f1d2a24d51de

## Where to see it?
Currently used by me as my personal blog, you can see it at https://mwyndham.dev/

## Features
All of the features include:
- WYSIWYG Editor w/ Markdown Keybinding (Support most Markdown formatting)
- Drag and Drop Image Upload
- Copy Paste Youtube Video Embed
- 20+ Theme (from DaisyUI) with Dark Mode Toggle
- RSS Feeds
- Sitemap Ping
- Hashtag
- S3 Support hosting resources (Image, GIF, JS, CSS)
- Auth (via Clerk)

It also has niceties for development:
- Auto hot reload (via Air)
- Client-side NPM Module support (via ESBuild)
- EZ Deployment to https://fly.io
- Beautiful LibSQL Database via https://turso.tech
- S3 Bucket hosted on Cloudflare R2

## Getting Started
Create a `.env` file (see `env.sample`), fill the ENV as required, then run `make local-setup` and `make run`. 

## Things to note
This blog is for my own personal use, hence many values are hardcoded to link to my own stuff (like links to Twitter, YouTube, Twitch, Medium, etc). Feel free to fork and modify it but don't say I don't warn you!
