import { ChatInputCommandInteraction } from "discord.js";

export function support(interaction: ChatInputCommandInteraction) {
    try {
        return interaction.reply({
            content: "Join my support server through this invite: https://discord.gg/S5yryjRuse",
            ephemeral: true
        });
    } catch { }
}