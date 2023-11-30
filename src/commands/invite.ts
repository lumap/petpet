import { ChatInputCommandInteraction } from "discord.js";

export function invite(interaction: ChatInputCommandInteraction) {
    try {
        return interaction.reply({
            content: "Click the button below to add me to one of your servers, or share this link to your friends to invite me: <https://a.lumap.cat /botinvite?id=966271016953327649&perms=51200>",
            components: [{
                type: 1,
                components: [
                    {
                        type: 2,
                        label: "Invite Me",
                        style: 5,
                        url: "https://discord.com/api/oauth2/authorize?client_id=966271016953327649&permissions=51200&scope=applications.commands%20bot"
                    }
                ]
            }],
            ephemeral: true
        });
    } catch { }
}