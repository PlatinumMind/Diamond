use serde::{Deserialize, Serialize};
use std::{
    env, fs,
    path::Path,
    process::{self, exit, Command},
};
use toml;

mod target;
use target::{c_bin, cpp_bin, workspace};

#[derive(Debug, Deserialize, Serialize)]
struct Target {
    cpp_bin: Option<Vec<cpp_bin::CxxBin>>,
    c_bin: Option<Vec<c_bin::CBin>>,
    workspace: Option<Vec<workspace::Workspace>>,
}

fn main() {
    let args: Vec<String> = env::args().collect();
    if args.len() != 1 {
        if args.len() == 3 && args[1] == "build" && !args[2].is_empty() {
            println!("Building workspace: {}", args[2]);
            build_with_workspace(args[2].to_string())
        } else if args.len() == 2 && args[1] == "build" {
            build();
        } else {
            println!("invalid arg");
            exit(1);
        }
    } else {
        println!("no arguments found");
        exit(1);
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
                    println!("Build compelted");
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
                    println!("Build compelted");
                }
            }
        }
        None => {}
    }
}

fn build_with_workspace(workspace_name: String) {
    let toml = run_diamond();
    match toml.workspace {
        Some(workspaces) => {
            for workspace in workspaces {
                if workspace.name == workspace_name {
                    let _ = env::set_current_dir(Path::new(&workspace.path));
                    build()
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
