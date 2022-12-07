import { CommandInteraction } from "discord.js";

export async function getUserPFP(interaction: CommandInteraction): Promise<string> {
    let content: string;
    if (interaction.guild) {
        try {
            const member = await interaction.guild.members.fetch(interaction.options.getUser("user")!);
            content = member.displayAvatarURL({ format: "png", size: 1024 });
        } catch {
            content = "This member doesn't seem to be here. If you want to petpet them, do it in my DMs.";
        }
    } else {
        content = interaction.options.getUser("user")!.avatarURL({ format: "png", size: 1024 })!;
    }
    return content;
}
