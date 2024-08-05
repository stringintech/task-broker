extern crate prost_build;

fn main() {
    prost_build::compile_protos(&["../proto/base/task.proto", "../proto/event/task_created.proto"], //TODO? explicit proto file names?
                                &["../proto"]).unwrap();
}