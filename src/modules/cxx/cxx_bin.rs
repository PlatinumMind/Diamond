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
        let cmd = gen_command(
            self.cxx,
            self.src,
            self.root,
            self.cxxflags,
            self.ldflags,
            self.output,
        );
        return cmd;
    }
}

fn gen_command(
    compiler: String,
    source: String,
    root: String,
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
        + &root
        + "/**/* "
        + &cxxflags
        + " "
        + &ldflags
        + "-o "
        + &output.out_dir
        + "/"
        + &output.exec_name;
}
