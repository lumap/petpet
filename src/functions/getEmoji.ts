import { CommandInteraction } from "discord.js";
const { parse } = require('twemoji-parser');

export function getEmoji(interaction: CommandInteraction) {
    let content;
    const emoji = interaction.options.getString("emoji")!;
    if ((emoji.match(/(<a?)?:\w+:(\d{16,20}>)?/u)) !== null) {
        content = `https://cdn.discordapp.com/emojis/${emoji.split(":")[2].slice(0, -1)}.png`;
    } else if (parse(emoji)?.[0]?.url) {
        content = parse(emoji)[0].url;
    } else {
        content = "I wasn't able to find an emoji in this. I wish discord had an \"emoji\" option for slash commands";
    }
    return content;

}
