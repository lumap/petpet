import { makePetGif } from "../functions/makePetGif";
import { isImage } from "../functions/isImage";
import { APIApplicationCommandInteraction, APIChatInputApplicationCommandInteractionData, ApplicationCommandOptionType } from "discord-api-types/v10";
import { FastifyReply } from "fastify";
import { deferInteration, editMessage, editMessageWithAttachments, sendMessage } from "../functions/interactions";
import { logCommand } from "../functions/logs";
let urlcheck = require('is-a-url');

export async function petpet(interaction: APIChatInputApplicationCommandInteractionData, res: FastifyReply, ogInteraction: APIApplicationCommandInteraction) {
    let param = interaction.options![0];
    if (!interaction.options || param.type !== ApplicationCommandOptionType.Subcommand || !param.options) return sendMessage(res, {
        content: "What the fuck?",
        flags: 64
    });
    let url = "", target = "";
    switch (param.options[0].type) {
        case ApplicationCommandOptionType.User: {
            const userId = param.options[0].value;
            const avatarHash = interaction.resolved?.members?.[userId].avatar || interaction.resolved?.users?.[userId].avatar;
            url = avatarHash ? `https://cdn.discordapp.com/avatars/${userId}/${avatarHash}.png?size=1024` : `https://cdn.discordapp.com/embed/avatars/${(Number(userId) >> 22) % 6}.png?size=1024`;
            target = interaction.resolved?.users?.[userId].username || "someone, idk discord fucked up";
            logCommand("petpet user");
            break;
        }
        case ApplicationCommandOptionType.Attachment: {
            const attachmentId = param.options[0].value;
            url = interaction.resolved?.attachments?.[attachmentId]!.url!;
            target = "an attachment";
            logCommand("petpet attachment");
            break;
        }
        case ApplicationCommandOptionType.String: {
            const str = param.options?.[0].value.split("?")[0];
            if (!urlcheck(str) || str.startsWith("https://tenor.com/view/") || !isImage(str)) {
                return sendMessage(res, {
                    content: "URL invalid",
                    flags: 64
                });
            }
            url = str;
            target = "an image url";
            logCommand("petpet imageurl");
            break;
        }
    }
    deferInteration(res);
    let gif: Buffer | string;
    try {
        let options = {
            resolution: param.options.find(c => c.type == ApplicationCommandOptionType.Number && c.name == "resolution")?.value || 128,
            delay: param.options.find(c => c.type == ApplicationCommandOptionType.Number && c.name == "delay")?.value || 30
        };
        gif = await makePetGif(url, options);
        if (typeof gif === "string") { throw new Error(gif); }
    } catch (e) {
        return editMessage(ogInteraction, {
            content: "Something fucked up, my bad. Please retry with something/someone else.",
        });
    }
    return editMessageWithAttachments(ogInteraction, {
        attachments: [
            {
                id: 0,
                filename: "pet.gif",
                description: `${ogInteraction.user?.username || "Someone"} has pet ${target}`
            }
        ]
    }, [gif]);
}