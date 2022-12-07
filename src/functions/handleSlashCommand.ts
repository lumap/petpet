import { CommandInteraction } from "discord.js";
import { invite } from "./invite";
import { sendGif } from "./sendGif";
import { getSlashURL } from "./getSlashURL";
import { createGif } from "./createGif";

export async function handleSlashCommand(interaction: CommandInteraction) {
    const ephemeral = interaction.options.getBoolean("ephemeral") || false;
    if (interaction.commandName === "invite") {
        invite(interaction);
        return;
    }
    await interaction.deferReply({ ephemeral: ephemeral });
    let content = await getSlashURL(interaction);
    if (!content.startsWith("http")) {
        return interaction.editReply({
            content: content
        })
    }
    let gif = await createGif(interaction, content);
    if (typeof gif === "boolean") return;
    sendGif(interaction, gif)
}