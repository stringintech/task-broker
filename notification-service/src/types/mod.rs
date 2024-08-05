pub mod event {
    include!(concat!(env!("OUT_DIR"), "/types.event.rs"));
}

pub mod base {
    include!(concat!(env!("OUT_DIR"), "/types.base.rs"));
}