use regex::Regex;
use std::{fmt, str::FromStr};

#[derive(Debug, Clone)]
pub struct PassageReference {
    pub book: String,
    pub start_chapter: i32,
    pub start_verse: i32,
    pub end_chapter: i32,
    pub end_verse: i32,
}
impl fmt::Display for PassageReference {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(
            f,
            "{} {}:{}-{}:{}",
            self.book, self.start_chapter, self.start_verse, self.end_chapter, self.end_verse
        )
    }
}

#[derive(Debug)]
pub struct PassageReferenceParseError;

impl FromStr for PassageReference {
    type Err = PassageReferenceParseError;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        let reference_regex: Regex = Regex::new(r"^((?:(?:1|2) )?\w+) (\d+):(\d+)-(\d+):(\d+)$")
            .expect("REFERENCE_REGEX failed to compile");
        let Some(reference) = reference_regex.captures(&s[..]).and_then(|captures| {
            let (Ok(start_chapter), Ok(start_verse), Ok(end_chapter), Ok(end_verse)) = (
                captures[2].parse(),
                captures[3].parse(),
                captures[4].parse(),
                captures[5].parse(),
            ) else {
                return None;
            };
            Some(PassageReference {
                book: String::from(&captures[1]),
                start_chapter,
                start_verse,
                end_chapter,
                end_verse,
            })
        }) else {
            return Err(PassageReferenceParseError);
        };
        Ok(reference)
    }
}

#[derive(Debug, Clone)]
pub struct Passage {
    pub id: i32,
    pub reference: PassageReference,
    pub level: i32,
}
