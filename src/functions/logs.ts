import axios from "axios";

export function logCommand(commandName: string) {
    axios(`${process.env.LOGS_WEBHOOK_URL}?thread_id=${process.env.LOGS_THREADS_COMMANDS}`, {
        method: "POST",
        data: {
            content: `Command **${commandName}** just got used.`
        },
        headers: {
            "Content-Type": "application/json"
        }
    });
}

export function logBoot() {
    axios(`${process.env.LOGS_WEBHOOK_URL}?thread_id=${process.env.LOGS_THREADS_BOOT}`, {
        method: "POST",
        data: {
            content: "Bot started!"
        },
        headers: {
            "Content-Type": "application/json"
        }
    });
}

export function logError(error: unknown) {
    axios(`${process.env.LOGS_WEBHOOK_URL}?thread_id=${process.env.LOGS_THREADS_ERRORS}`, {
        method: "POST",
        data: {
            content: `**${error || "An error happened"}**`
        },
        headers: {
            "Content-Type": "application/json"
        }
    });
}