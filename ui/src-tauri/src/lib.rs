use tauri::{WebviewWindow, Manager};
use std::env;

#[tauri::command]
async fn open_url(url: String) -> Result<(), String> {
    tauri_plugin_opener::open_url(url, None::<&str>).map_err(|e| e.to_string())
}

#[tauri::command]
async fn hide_window(window: WebviewWindow) -> Result<(), String> {
    window.hide().map_err(|e| e.to_string())
}

fn parse_shortcut(shortcut_str: &str) -> Option<(Option<tauri_plugin_global_shortcut::Modifiers>, tauri_plugin_global_shortcut::Code)> {
    let parts: Vec<&str> = shortcut_str.to_lowercase().split('+').collect();
    if parts.len() != 2 {
        return None;
    }
    
    let modifier = match parts[0] {
        "cmd" => Some(tauri_plugin_global_shortcut::Modifiers::META),
        "ctrl" => Some(tauri_plugin_global_shortcut::Modifiers::CONTROL),
        "alt" => Some(tauri_plugin_global_shortcut::Modifiers::ALT),
        "shift" => Some(tauri_plugin_global_shortcut::Modifiers::SHIFT),
        _ => return None,
    };
    
    let code = match parts[1] {
        "space" => tauri_plugin_global_shortcut::Code::Space,
        _ => return None,
    };
    
    Some((modifier, code))
}

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    // Load environment variables
    dotenv::dotenv().ok();
    tauri::Builder::default()
        .plugin(tauri_plugin_opener::init())
        .invoke_handler(tauri::generate_handler![open_url, hide_window])
        .setup(|app| {
            #[cfg(desktop)]
            {
                use tauri_plugin_global_shortcut::{Code, GlobalShortcutExt, Modifiers, Shortcut, ShortcutState};

                // Get shortcut from environment variable or use default
                let shortcut_str = env::var("SHORTCUT").unwrap_or_else(|_| {
                    #[cfg(target_os = "macos")]
                    return "cmd+space".to_string();
                    #[cfg(not(target_os = "macos"))]
                    return "ctrl+space".to_string();
                });

                let shortcut = if let Some((modifier, code)) = parse_shortcut(&shortcut_str) {
                    Shortcut::new(modifier, code)
                } else {
                    // Fallback to default if parsing fails
                    #[cfg(target_os = "macos")]
                    let fallback = Shortcut::new(Some(Modifiers::META), Code::Space);
                    #[cfg(not(target_os = "macos"))]
                    let fallback = Shortcut::new(Some(Modifiers::CONTROL), Code::Space);
                    fallback
                };

                let app_handle = app.handle().clone();

                // Set up window event listener to hide when focus is lost
                if let Some(window) = app.get_webview_window("main") {
                    let app_handle_for_focus = app_handle.clone(); // Clone for the focus handler
                    window.on_window_event(move |event| {
                        match event {
                            tauri::WindowEvent::Focused(focused) => {
                                if !focused {
                                    // Window lost focus, hide it
                                    if let Some(window) = app_handle_for_focus.get_webview_window("main") {
                                        let _ = window.hide();
                                    }
                                }
                            }
                            _ => {}
                        }
                    });
                }

                app.handle().plugin(
                    tauri_plugin_global_shortcut::Builder::new()
                        .with_handler(move |_app, pressed_shortcut, event| {
                            if pressed_shortcut == &shortcut && event.state() == ShortcutState::Pressed {
                                if let Some(window) = app_handle.get_webview_window("main") {
                                    let _ = window.show();
                                    let _ = window.set_focus();
                                }
                            }
                        })
                        .build(),
                )?;

                app.global_shortcut().register(shortcut)?;
            }

            Ok(())
        })
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
