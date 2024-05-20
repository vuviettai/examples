fn main() {
    println!("cargo:rerun-if-changed=proto");
    //tonic_build::compile_protos("../proto/message.proto").unwrap();
    tonic_build::configure()
        .out_dir("src/proto")
        .compile(&["../proto/message.proto"], &["../proto"])
        .expect("Failed to compile proto(s)");
}
