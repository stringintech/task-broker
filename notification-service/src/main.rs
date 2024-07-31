use prost::Message;
use lapin::{options::*, types::FieldTable, Connection, ConnectionProperties};
use futures_util::stream::StreamExt;
use std::error::Error;
use std::io::Cursor;

mod types {
    include!(concat!(env!("OUT_DIR"), "/types.rs"));
}

use types::Task;

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

    println!("waiting for task notifications");

    while let Some(delivery) = consumer.next().await {
        match delivery {
            Ok(delivery) => {
                let task = Task::decode(&mut Cursor::new(&delivery.data))?;
                println!("received task notification");
                println!("task id: {}", task.id);
                println!("task title: {}", task.title);
                delivery.ack(BasicAckOptions::default()).await?;
            }
            Err(error) => {
                eprintln!("error receiving task notification: {:?}", error);
            }
        }
    }

    Ok(())
}
