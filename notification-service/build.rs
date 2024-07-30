extern crate prost_build;

fn main() {
    prost_build::compile_protos(&["../proto/message.proto"], //TODO? explicit proto file names?
                                &["../proto"]).unwrap();
}