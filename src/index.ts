import { Client, GatewayIntentBits } from "discord.js";
import commands from "./commands";
const config = require("../config.js");
const client = new Client({ intents: [GatewayIntentBits.Guilds], ws: { properties: { browser: "Discord iOS" } } });

client.on("interactionCreate", async (interaction) => {
    if (!interaction.isCommand()) return;
    if (interaction.isUserContextMenuCommand()) {
        commands.petpetUserCtx(interaction, client);
    } else if (interaction.isMessageContextMenuCommand()) {
        commands.petpetMsgCtx(interaction, client);
    } else {
        switch (interaction.commandName) {
            case "invite": {
                commands.invite(interaction);
                break;
            }
            case "support": {
                commands.support(interaction);
                break;
            }
            case "petpet": {
                commands.petpet(interaction);
                break;
            }
        }
    }
});

client.login(config.token);
console.log("Starting...");
client.on("ready", () => {
    console.log("bot started ig");
});