use crate::Focus;

#[derive(Debug, Default)]
pub struct InputBox {
    pub text: String,
}

impl InputBox {
    pub fn own_focus(&self) -> Focus {
        Focus::InputBox
    }

    pub fn next_focus(&self) -> Focus {
        Focus::MovieTable
    }
}
