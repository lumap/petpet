import { makePetGif } from "../functions/makePetGif";
import { APIApplicationCommandInteraction, APIMessageApplicationCommandInteractionData } from "discord-api-types/v10";
import { FastifyReply } from "fastify";
import { deferInteration, editMessage, editMessageWithAttachments } from "../functions/interactions";
import { logCommand } from "../functions/logs";

export async function petpetMsgCtx(interaction: APIMessageApplicationCommandInteractionData, res: FastifyReply, ogInteraction: APIApplicationCommandInteraction) {
    logCommand("petpetMsgCtx");
    const msg = interaction.resolved.messages[interaction.target_id];
    const avatarHash = msg.author.avatar;
    const url = avatarHash ? `https://cdn.discordapp.com/avatars/${msg.author.id}/${avatarHash}.png?size=1024` : `https://cdn.discordapp.com/embed/avatars/${(Number(msg.author.id) >> 22) % 6}.png?size=1024`;
    await deferInteration(res);
    let gif: Buffer | string;
    try {
        let options = {
            resolution: 128,
            delay: 30
        };
        gif = await makePetGif(url, options);
        if (typeof gif === "string") { return; }
    } catch {
        return editMessage(ogInteraction, {
            content: "Something fucked up, my bad. Please retry with something/someone else.",
        });
    }
    return editMessageWithAttachments(ogInteraction, {
        attachments: [
            {
                id: 0,
                filename: "pet.gif",
                description: `${ogInteraction.member?.user.username || ogInteraction.user?.username || "Someone"} has pet ${msg.author.username}`
            }
        ]
    }, [gif]);
}