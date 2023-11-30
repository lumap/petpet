import { ChatInputCommandInteraction } from "discord.js";
import { makePetGif } from "../functions/makePetGif";
import { isImage } from "../functions/isImage";
import { sendGif } from "../functions/sendGif";
let urlcheck = require('is-a-url');

async function getSlashURL(interaction: ChatInputCommandInteraction): Promise<{ content: string, target: string; }> {
    let content = "", target = "";
    switch (interaction.options.getSubcommand()) {
        case "user": {
            const user = interaction.options.getUser("user")!;
            if (interaction.guild) {
                try {
                    const member = await interaction.guild.members.fetch(user);
                    target = member.user.tag;
                    content = member.displayAvatarURL({ extension: "png", size: 1024 });
                } catch {
                    target = user.tag;
                    content = user.displayAvatarURL({ extension: "png", size: 1024 })!;
                }
            } else {
                target = user.tag;
                content = user.displayAvatarURL({ extension: "png", size: 1024 })!;
            }
            break;
        }
        case "attachment": {
            const url = interaction.options.getAttachment("attachment")!.url.split("?")[0];
            if (!urlcheck(url) || !isImage(url)) {
                content = "Sorry, this attachment does not seem to be valid. Please make sure it's a `jpg`, `jpeg` or `png` image.";
            } else {
                content = url;
            }
            target = "an attachment";
            break;
        }
        case "imageurl": {
            const url = interaction.options.getString("url")!;
            if (!urlcheck(url) || url.startsWith("https://tenor.com/view/") || !isImage(url)) {
                content = "Sorry, this link does not seem to be valid. Please make sure the image link ends with `.jpg`, `.jpeg` or `.png`.";
            } else {
                content = url;
            }
            target = "an image from an external URL";
            break;
        }
        case "server": {
            if (!interaction.guild) {
                content = "This command is only possible in servers";
                break;
            }
            const url = interaction.guild.iconURL({ extension: "png", size: 1024 });
            if (!url) {
                content = "This server does't have an icon.";
                break;
            }
            content = url;
            target = "this server's icon";
            break;
        }
        default: {
            break;
        }
    }
    return { content, target };
}

export async function petpet(interaction: ChatInputCommandInteraction) {
    const ephemeral = interaction.options.getBoolean("ephemeral") || false;
    await interaction.deferReply({ ephemeral: ephemeral });
    let { content, target } = await getSlashURL(interaction);
    if (!content?.startsWith("http")) {
        return interaction.editReply({
            content: content
        });
    }
    let gif: Buffer | string;
    try {
        let options = {
            resolution: 128,
            delay: 30
        };
        if (interaction.isChatInputCommand()) {
            if (interaction.options.getInteger("delay")) {
                options.delay = interaction.options.getInteger("delay")!;
            }
            if (interaction.options.getInteger("resolution")) {
                options.resolution = interaction.options.getInteger("resolution")!;
            }
        }
        gif = await makePetGif(content, options);
        if (typeof gif === "string") { return; }
    } catch {
        interaction.editReply({
            content: "Sorry, but it looks like something went wrong. Please retry with a valid file/link",
        });
        return;
    }
    sendGif(interaction, gif!, target);
}