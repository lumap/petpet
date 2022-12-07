import { Client, UserContextMenuInteraction } from "discord.js";
import { getPetGif } from "./getPetGif";
import { sendGif } from "./sendGif"

export async function handleUserContextMenu(interaction: UserContextMenuInteraction, client: Client) {
    await interaction.deferReply();
    let content: string;
    if (interaction.guild) {
        try {
            const member = await interaction.guild.members.fetch(interaction.targetId);
            content = member.displayAvatarURL({ format: "png", size: 1024 });
        } catch {
            content = "This member doesn't seem to be here. If you want to petpet them, do it in my DMs.";
        }
    } else {
        const user = await client.users.fetch(interaction.targetId);
        content = user.avatarURL({ format: "png", size: 1024 })!;
    }
    if (!content.startsWith("http")) {
        return interaction.editReply({
            content: content
        });
    }
    const gif = await getPetGif(content);
    if (typeof gif === "string") {
        return interaction.editReply({ content: "I fucked up" });
    }
    sendGif(interaction, gif)
}
