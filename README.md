# wistia-dl
You can use this tool to download videos hosted on the [Wistia](https://wistia.com/) platform.

## Download
You can find the binaries in the [Releases](https://github.com/smbl64/wistia-dl/releases) page. Windows, macOS and Linux are supported.
## Usage
The command is quite simple: `wistia-dl -v <video_id> -o <output_file_name>`.

To find the video ID, follow these steps:

1. Right click on the video and choose "Copy link and thumbnail".
![Copy link and thumbnail](./docs/click-on-video.png)

2. Paste the copied text in a text editor and find the video ID. Look for `"wvideo="`. Video ID comes after that text. In the following picture, the video ID is `0j69dfdsq4`.
![Video ID](./docs/video-id.png)


## Build from Source
1. Install Go.
2. Clone the project.
3. Run the `build.sh`. It will generate binary files for macOS, Linux, and Windows in the `output` folder.
