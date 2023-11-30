import { Client, UserContextMenuCommandInteraction } from "discord.js";
import { makePetGif } from "../functions/makePetGif";
import { sendGif } from "../functions/sendGif";


export async function petpetUserCtx(interaction: UserContextMenuCommandInteraction, client: Client) {
    try { await interaction.deferReply(); } catch { return; }
    let content: string = "nice try";
    let target: string = "h";
    if (interaction.guild) {
        try {
            const member = await interaction.guild.members.fetch(interaction.targetId);
            target = member.user.tag;
            content = member.displayAvatarURL({ extension: "png", size: 1024 });
        } catch {
            content = "This member doesn't seem to be here. If you want to petpet them, use their user ID as the `user` argument of `/petpet user`.";
        }
    } else {
        const user = await client.users.fetch(interaction.targetId);
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
        return interaction.editReply({
            content: "I fucked up"
        });
    }
    sendGif(interaction, gif, target);
}