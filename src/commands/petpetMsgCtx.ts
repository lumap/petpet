import { Client, MessageContextMenuCommandInteraction } from "discord.js";
import { sendGif } from "../functions/sendGif";
import { makePetGif } from "../functions/makePetGif";

export async function petpetMsgCtx(interaction: MessageContextMenuCommandInteraction, client: Client) {
    try { await interaction.deferReply(); } catch { "why do u keep crashing here"; return; }
    let content: string = "nice try", target: string;
    if (interaction.guild) {
        try {
            if (interaction.targetMessage.webhookId) {
                target = interaction.targetMessage.author.username + " (webhook)";
                content = interaction.targetMessage.author.displayAvatarURL({ extension: "png", size: 1024 });
            } else {
                const member = await interaction.guild.members.fetch(interaction.targetMessage.author.id);
                target = member?.user?.tag;
                content = member.displayAvatarURL({ extension: "png", size: 1024 });
            }
        } catch (e: any) {
            console.log(e);
            content = "This member doesn't seem to be here. If you want to petpet them, use their user ID as the `user` argument of `/petpet user`.";
            target = "h";
        }
    } else {
        const user = await client.users.fetch(interaction.targetMessage.author.id);
        target = user.tag;
        content = user.avatarURL({ extension: "png", size: 1024 })!;
    }
    if (!content.startsWith("http")) {
        return interaction.editReply({
            content: content
        });
    }
    let options = {
        resolution: 128,
        delay: 30
    };
    const gif = await makePetGif(content, options);
    if (typeof gif === "string") {
        return interaction.editReply({ content: "I fucked up" });
    }
    sendGif(interaction, gif, target);
}