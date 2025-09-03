use tauri::{WebviewWindow, Manager};

#[tauri::command]
async fn open_url(url: String) -> Result<(), String> {
    tauri_plugin_opener::open_url(url, None::<&str>).map_err(|e| e.to_string())
}

#[tauri::command]
async fn hide_window(window: WebviewWindow) -> Result<(), String> {
    window.hide().map_err(|e| e.to_string())
}

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    tauri::Builder::default()
        .plugin(tauri_plugin_opener::init())
        .invoke_handler(tauri::generate_handler![open_url, hide_window])
        .setup(|app| {
            #[cfg(desktop)]
            {
                use tauri_plugin_global_shortcut::{Code, GlobalShortcutExt, Modifiers, Shortcut, ShortcutState};

                // Use Cmd on macOS, Ctrl on Windows/Linux
                #[cfg(target_os = "macos")]
                let shortcut = Shortcut::new(Some(Modifiers::META), Code::Space);

                #[cfg(not(target_os = "macos"))]
                let shortcut = Shortcut::new(Some(Modifiers::CONTROL), Code::Space);

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
