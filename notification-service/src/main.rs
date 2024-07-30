use lapin::{options::*, types::FieldTable, Connection, ConnectionProperties};
use futures_util::stream::StreamExt;
use std::error::Error;

#[tokio::main]
async fn main() -> Result<(), Box<dyn Error>> {
    let addr = "amqp://guest:guest@localhost:5672//";
    let conn = Connection::connect(&addr, ConnectionProperties::default()).await?;
    let channel = conn.create_channel().await?;

    let _queue = channel.queue_declare(
        "hello",
        QueueDeclareOptions::default(),
        FieldTable::default(),
    ).await?;

    let mut consumer = channel.basic_consume(
        "hello",
        "my_consumer",
        BasicConsumeOptions::default(),
        FieldTable::default(),
    ).await?;

    println!("waiting for messages");

    while let Some(delivery) = consumer.next().await {
        match delivery {
            Ok(delivery) => {
                println!("received message: {}", String::from_utf8_lossy(&delivery.data));
                delivery.ack(BasicAckOptions::default()).await?;
            }
            Err(error) => {
                eprintln!("error receiving message: {:?}", error);
            }
        }
    }

    Ok(())
}
