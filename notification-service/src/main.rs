mod types;
mod services;

use prost::Message;
use lapin::{options::*, types::FieldTable, Connection, ConnectionProperties};
use futures_util::stream::StreamExt;
use std::error::Error;
use std::io::Cursor;
use types::event;
use crate::services::event_handler::EventHandler;
use crate::services::notification_service::EmailNotificationService;

#[tokio::main]
async fn main() -> Result<(), Box<dyn Error>> {
    let addr = "amqp://guest:guest@localhost:5672//";
    let conn = Connection::connect(&addr, ConnectionProperties::default()).await?;
    let channel = conn.create_channel().await?;

    let _ = channel.queue_declare(
        "task-queue",
        QueueDeclareOptions::default(),
        FieldTable::default(),
    ).await?;

    let mut consumer = channel.basic_consume(
        "task-queue",
        "dummy-consumer",
        BasicConsumeOptions::default(),
        FieldTable::default(),
    ).await?;

    println!("waiting for task events");

    let email_service = EmailNotificationService::new("localhost".parse().unwrap(), 1025);
    let handler = EventHandler::new(Box::new(email_service));

    // Consume messages
    while let Some(delivery) = consumer.next().await {
        match delivery {
            Ok(delivery) => {
                let event = event::TaskCreated::decode(&mut Cursor::new(&delivery.data))?;
                match handler.handle_event_task_created(&event) {
                    Ok(()) => {
                        println!("event handled");
                        delivery.ack(BasicAckOptions::default()).await?;
                    }
                    Err(error) => {
                        eprintln!("error handling event: {:?}", error);
                    }
                }
            }
            Err(error) => {
                eprintln!("error receiving task event: {:?}", error);
            }
        }
    }

    Ok(())
}
