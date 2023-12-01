import { FastifyReply } from "fastify";
import { sendMessage } from "../functions/interactions";
import { APIChatInputApplicationCommandInteractionData, ButtonStyle, ComponentType } from "discord-api-types/v10";
import { logCommand } from "../functions/logs";

export function invite(interaction: APIChatInputApplicationCommandInteractionData, res: FastifyReply) {
    logCommand("invite");
    return sendMessage(res, {
        content: "Click the button below to add me to one of your servers, or share this link to your friends to invite me: <https://a.lumap.cat /botinvite?id=966271016953327649&perms=51200>",
        components: [{
            type: ComponentType.ActionRow,
            components: [
                {
                    type: ComponentType.Button,
                    label: "Invite Me",
                    style: ButtonStyle.Link,
                    url: "https://discord.com/api/oauth2/authorize?client_id=966271016953327649&permissions=51200&scope=applications.commands%20bot"
                }
            ]
        }],
        flags: 64
    });
}