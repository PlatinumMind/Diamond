use std::fs;
use serde::{Deserialize, Serialize};
use toml;

mod modules;
use modules::cxx_bin;

#[derive(Debug, Deserialize, Serialize)]
struct Target {
    cpp_bin: Vec<cxx_bin::CxxBin>
}

fn main() {
    let binding = fs::read_to_string("./diamond.build").unwrap();
    let raw_toml = binding.as_str();
    let toml: Target = toml::from_str(raw_toml).unwrap();
    println!("{:?}", toml.cpp_bin);
    for c in toml.cpp_bin {
        println!("{}", c.make_command())
    }
}
