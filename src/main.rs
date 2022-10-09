use serde::{Deserialize, Serialize};
use std::{
    fs,
    process::{self, Command}, env,
};
use toml;

mod target;
use target::{c_bin, cpp_bin};

#[derive(Debug, Deserialize, Serialize)]
struct Target {
    cpp_bin: Option<Vec<cpp_bin::CxxBin>>,
    c_bin: Option<Vec<c_bin::CBin>>,
}

fn main() {
    let args = env::args();
    for arg in args {
        if arg == "build" {
            build();
        }
    }
}

fn build() {
    let toml = run_diamond();
    match toml.cpp_bin {
        Some(cpp_) => {
            for cpp in cpp_ {
                let cmd = format!("{}", cpp.make_command());
                let output = Command::new("bash")
                    .args(["-c", &cmd])
                    .output()
                    .expect("failed to build");

                if !output.stderr.is_empty() {
                    eprintln!("{}", String::from_utf8_lossy(&output.stderr));
                    process::exit(1);
                } else {
                    print!("{}", String::from_utf8_lossy(&output.stdout));
                    println!(" Build compelted");

                }
            }
        }
        None => {}
    }
    match toml.c_bin {
        Some(c_) => {
            for c in c_ {
                let cmd = format!("{}", c.make_command());
                let output = Command::new("bash")
                    .args(["-c", &cmd])
                    .output()
                    .expect("failed to build");

                if !output.stderr.is_empty() {
                    eprintln!("{}", String::from_utf8_lossy(&output.stderr));
                    process::exit(1);
                } else {
                    print!("{}", String::from_utf8_lossy(&output.stdout));
                    println!(" Build compelted");
                }
            }
        }
        None => {}
    }
}


fn run_diamond() -> Target {
    let binding = fs::read_to_string("./diamond.build").unwrap();
    let raw_toml = binding.as_str();
    let toml: Target = toml::from_str(raw_toml).unwrap();
    return toml;
}
