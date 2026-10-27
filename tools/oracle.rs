// SPDX-License-Identifier: GPL-3.0-or-later
// Go port and modifications of chaptr 1.5.0, 2026-10-27, by emirrcaglar.
// Upstream: https://github.com/johnthreekay/chaptr. See NOTICE for attribution.

// Test fixture generator only. The Go library has no Rust runtime dependency.
use std::io::{self, Read};
use chaptr::{ChapterNumber, NumberRange};

fn string(s: &str) -> String {
    let mut out = String::from("\"");
    for c in s.chars() {
        match c {
            '"' => out.push_str("\\\""),
            '\\' => out.push_str("\\\\"),
            '\n' => out.push_str("\\n"),
            '\r' => out.push_str("\\r"),
            '\t' => out.push_str("\\t"),
            c if c < '\u{20}' => out.push_str(&format!("\\u{:04x}", c as u32)),
            c => out.push(c),
        }
    }
    out.push('"'); out
}
fn text(s: Option<&str>) -> String { string(s.unwrap_or("")) }
fn option<T: std::fmt::Display>(n: Option<T>) -> String { n.map(|n| n.to_string()).unwrap_or("null".into()) }
fn number(n: ChapterNumber) -> String { format!("{{\"whole\":{},\"decimal\":{}}}",n.whole,option(n.decimal)) }
fn range(r: Option<NumberRange>) -> String {
    r.map(|r| format!("{{\"start\":{},\"end\":{}}}",number(r.start),r.end.map(number).unwrap_or("null".into()))).unwrap_or("null".into())
}
fn main() {
    let mut input = String::new(); io::stdin().read_to_string(&mut input).unwrap();
    for s in input.split('\0') {
        let m=chaptr::manga::parse(s); let n=chaptr::novel::parse(s);
        let manga=format!("{{\"title\":{},\"volume\":{},\"chapter\":{},\"group\":{},\"source\":{},\"language\":{},\"revision\":{},\"is_oneshot\":{},\"edition\":{},\"extension\":{}}}",
            text(m.title.as_deref()),range(m.volume),range(m.chapter),text(m.group),string(&m.source.map(|x|format!("{x:?}")).unwrap_or_default()),string(&m.language.map(|x|format!("{x:?}")).unwrap_or_default()),option(m.revision),m.is_oneshot,text(m.edition.as_deref()),text(m.extension));
        let novel=format!("{{\"title\":{},\"volume\":{},\"group\":{},\"publisher\":{},\"scanner\":{},\"language\":{},\"year\":{},\"is_premium\":{},\"is_digital\":{},\"revision\":{},\"extension\":{}}}",
            text(n.title.as_deref()),range(n.volume),text(n.group),text(n.publisher),text(n.scanner),string(&n.language.map(|x|format!("{x:?}")).unwrap_or_default()),option(n.year),n.is_premium,n.is_digital,option(n.revision),text(n.extension));
        let tokens=chaptr::lexer::tokenize(s).iter().map(|t| {
            use chaptr::lexer::Token::*;
            let (kind,txt,start,end)=match t {
                Word(w)=>(0,w.to_string(),w.as_ptr() as usize-s.as_ptr() as usize,w.as_ptr() as usize-s.as_ptr() as usize+w.len()),
                Delimiter(c)=>(1,c.to_string(),0,0),
                Bracketed(w)|Parenthesized(w)|Curly(w)=>{
                    let start=w.as_ptr() as usize-s.as_ptr() as usize-1;
                    (match t {Bracketed(_)=>2,Parenthesized(_)=>3,_=>4},w.to_string(),start,(start+w.len()+2).min(s.len()))
                }
            }; (kind,txt,start,end)
        }).collect::<Vec<_>>();
        let mut cursor=0;
        let tokens=tokens.iter().map(|(kind,txt,start,end)| {
            let (start,end)=if *kind==1 {(cursor,cursor+txt.len())} else {(*start,*end)};
            cursor=end;
            format!("{{\"kind\":{kind},\"text\":{},\"start\":{start},\"end\":{end}}}",string(txt))
        }).collect::<Vec<_>>().join(",");
        println!("{{\"input\":{},\"manga\":{manga},\"novel\":{novel},\"tokens\":[{tokens}]}}",string(s));
    }
}
