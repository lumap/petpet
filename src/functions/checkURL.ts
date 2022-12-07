import { CommandInteraction } from "discord.js";
import { isImage } from "./isImage";
let urlcheck = require('is-a-url');

export function checkURL(interaction: CommandInteraction): string {
    let content;
    const url = interaction.options.getString("url")!;
    if (!urlcheck(url) || url.startsWith("https://tenor.com/view/") || !isImage(url)) {
        content = "Sorry, this link does not seem to be valid. Please make sure the image link ends with `.jpg`, `.jpeg` or `.png`.";
    } else {
        content = url;
    }
    return content;
}
