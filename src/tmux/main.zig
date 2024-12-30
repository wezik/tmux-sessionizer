const std = @import("std");
const config = @import("config.zig");
const domain = @import("domain.zig");
pub const TmuxSession = domain.TmuxSession;
pub const TmuxEnvironment = domain.TmuxEnvironment;
pub const TmuxWindow = domain.TmuxWindow;

pub fn getSessions(allocator: std.mem.Allocator) ![]const TmuxSession {
    return try config.fetchConfig(allocator);
}

pub fn getConfigPath(allocator: std.mem.Allocator) ![]const u8 {
    return try config.getFullConfigPath(allocator);
}

pub fn appendSession(allocator: std.mem.Allocator, session: TmuxSession) !void {
    const sessions = try getSessions(allocator);
    var new_sessions = std.ArrayList(TmuxSession).init(allocator);
    for (sessions) |entry| {
        if (std.mem.eql(u8, entry.name, session.name)) {
            std.debug.print("Already exists: {s}\n", .{entry.name});
            return;
        }
        _ = try new_sessions.append(entry);
    }
    _ = try new_sessions.append(session);
    try config.saveConfig(allocator, try new_sessions.toOwnedSlice());
}

pub fn prepareSession(session: TmuxSession) !void {
    var arena = std.heap.ArenaAllocator.init(std.heap.page_allocator);
    defer arena.deinit();
    const allocator = arena.allocator();

    // check if session exists, if not create it
    var exists_cmd = std.process.Child.init(&[_][]const u8{ "tmux", "has-session", "-t", session.name }, allocator);
    exists_cmd.stderr_behavior = .Ignore;
    _ = try exists_cmd.spawn();
    const exit_code = try exists_cmd.wait();
    if (exit_code.Exited == 0) {
        // session exists, no need to prepare anything
        return;
    }

    _ = try helper_initSession(allocator, session);
    return;
}

fn helper_initSession(allocator: std.mem.Allocator, session: TmuxSession) !void {
    const first_window = session.windows[0];
    // create session
    const tmux_session_cmd = &[_][]const u8{
        "tmux",
        "new-session",
        "-d",
        "-s",
        session.name,
        "-c",
        session.path,
        "-n",
        first_window.name,
    };
    var tmux_session_exec = std.process.Child.init(tmux_session_cmd, allocator);
    _ = try tmux_session_exec.spawn();
    // wait for tmux to finish
    _ = try tmux_session_exec.wait();
    // open windows
    _ = try helper_openWindows(allocator, session);
    // set env in windows
    _ = try setEnvInWindows(allocator, session);
}

fn helper_openWindows(allocator: std.mem.Allocator, session: TmuxSession) !void {
    // dont open the first one as its the default one in the session
    for (session.windows[1..]) |window| {
        const new_window_cmd = &[_][]const u8{
            "tmux",
            "new-window",
            "-t",
            session.name,
            "-c",
            session.path,
            "-n",
            window.name,
        };
        var new_window_exec = std.process.Child.init(new_window_cmd, allocator);
        _ = try new_window_exec.spawn();
        // wait for tmux to finish
        _ = try new_window_exec.wait();
    }
}

fn setEnvInWindows(allocator: std.mem.Allocator, session: TmuxSession) !void {
    for (session.windows) |window| {
        const window_id = try std.fmt.allocPrint(allocator, "{s}:{s}", .{ session.name, window.name });

        // apply per session env
        for (session.env) |cmd| {
            _ = try helper_execCmdInWindow(allocator, cmd, window_id);
        }
        // apply per window env
        for (window.env) |cmd| {
            _ = try helper_execCmdInWindow(allocator, cmd, window_id);
        }
    }
}

fn helper_execCmdInWindow(allocator: std.mem.Allocator, cmd: []const u8, window_id: []const u8) !void {
    const in_window_cmd = &[_][]const u8{
        "tmux",
        "send-keys",
        "-t",
        window_id,
        cmd,
        "C-m",
    };
    var in_window_exec = std.process.Child.init(in_window_cmd, allocator);
    _ = try in_window_exec.spawn();
    // wait for tmux to finish
    _ = try in_window_exec.wait();
}

pub fn deleteSession(key: []const u8, allocator: std.mem.Allocator) !void {
    const sessions = try getSessions(allocator);
    var new_sessions = std.ArrayList(TmuxSession).init(allocator);
    for (sessions) |entry| {
        if (!std.mem.eql(u8, try sessionToKey(allocator, entry), key)) {
            _ = try new_sessions.append(entry);
        } else {
            std.debug.print("Removing: {s}\n", .{entry.name});
        }
    }
    _ = try config.saveConfig(allocator, try new_sessions.toOwnedSlice());
}

pub fn sessionToKey(allocator: std.mem.Allocator, entry: TmuxSession) ![]const u8 {
    if (std.mem.eql(u8, entry.name, entry.path)) {
        return entry.name;
    } else {
        const result = try std.fmt.allocPrint(allocator, "{s} ({s})", .{ entry.name, entry.path });
        return result;
    }
}
