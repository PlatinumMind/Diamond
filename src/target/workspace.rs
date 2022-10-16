use serde::{Deserialize, Serialize};

#[derive(Deserialize, Serialize, Debug)]
pub struct Workspace {
    pub name: String,
    pub path: String,
}   
