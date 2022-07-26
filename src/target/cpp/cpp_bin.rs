use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize)]
pub struct Output {
    exec_name: String,
    out_dir: String,
}

#[derive(Debug, Deserialize, Serialize)]
#[warn(non_camel_case_types)]
pub struct CxxBin {
    pub src: String,
    pub root: String,
    pub cxx: String,
    pub cache: Option<bool>,
    pub cxxflags: Option<Vec<String>>,
    pub ldflags: Option<Vec<String>>,
    pub output: Output,
}

impl CxxBin {
    pub fn make_command(self) -> String {
        let path = format!("{}{}", self.root.as_str(), "/**/*.cc");
        let build_path = self.output.out_dir.as_str();
        let mut files: Vec<String> = vec![];
        if !std::path::Path::new(&build_path).exists() {
            let _res = std::fs::create_dir(build_path);
        }
        for file in glob::glob(&path).expect("Failed to read glob pattern") {
            match file {
                Ok(path) => {
                    let path_str = path.display().to_string();
                    let local_src = self.src.strip_prefix("./").unwrap();
                    if path_str != local_src {
                        files.push(path.display().to_string())
                    }
                }
                Err(e) => println!("{:?}", e),
            }
        }
        let cmd = gen_command(
            self.cxx,
            self.src,
            files,
            self.cxxflags,
            self.ldflags,
            self.output,
        );
        if self.cache.unwrap() {
            return format!("CCACHE_DIR=.dbuild ccache {}", cmd)
        }
        return cmd;
    }
}

fn gen_command(
    compiler: String,
    source: String,
    files: Vec<String>,
    raw_cxxflags: Option<Vec<String>>,
    raw_ldflags: Option<Vec<String>>,
    output: Output,
) -> String {
    let cxxflags;
    let ldflags;
    match raw_cxxflags {
        Some(e) => {
            cxxflags = e.join(" ");
        }
        None => {
            cxxflags = "".to_string();
        }
    }
    match raw_ldflags {
        Some(e) => {
            ldflags = e.join(" ");
        }
        None => {
            ldflags = "".to_string();
        }
    }

    return compiler
        + " "
        + &source
        + " "
        + &files.join(" ")
        + " "
        + &cxxflags
        + " "
        + &ldflags
        + " -o "
        + &output.out_dir
        + "/"
        + &output.exec_name;
}
