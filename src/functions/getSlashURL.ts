import { CommandInteraction } from "discord.js";
import { getUserPFP } from "./getUserPFP";
import { getAttachment } from "./getAttachment";
import { checkURL } from "./checkURL";
import { getEmoji } from "./getEmoji";

export async function getSlashURL(interaction: CommandInteraction): Promise<string> {
    let content: string = "tf did u do";
    switch (interaction.options.getSubcommand()) {
        case "user": {
            content = await getUserPFP(interaction);
            break;
        }
        case "attachment": {
            content = getAttachment(interaction);
        }
        case "imageurl": {
            content = checkURL(interaction);
        }
        case "emoji": {
            content = getEmoji(interaction);
        }
        default: {
            break;
        }
    }
    return content;
}
