import { CommandInteraction } from "discord.js";
import { isImage } from "./isImage";
let urlcheck = require('is-a-url');

export function getAttachment(interaction: CommandInteraction): string {
    let content: string;
    const url = interaction.options.getAttachment("attachment")!.url;
    if (!urlcheck(url) || !isImage(url)) {
        content = "Sorry, this attachment does not seem to be valid. Please make sure it's a `jpg`, `jpeg` or `png` image.";
    } else {
        content = url;
    }
    return content;
}
