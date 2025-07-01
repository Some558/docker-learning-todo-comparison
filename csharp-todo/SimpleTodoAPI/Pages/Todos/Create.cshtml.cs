using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Mvc.RazorPages;
using SimpleTodoAPI.Data;
using SimpleTodoAPI.Models;

namespace SimpleTodoAPI.Pages.Todos
{
    public class CreateModel : PageModel
    {
        private readonly TodoContext _context;

        public CreateModel(TodoContext context)
        {
            _context = context;
        }

        public IActionResult OnGet()
        {
            Todo = new Todo(); // モデルを初期化
            return Page();
        }

        [BindProperty]
        public Todo Todo { get; set; } = new Todo(); // デフォルト値を設定

        public async Task<IActionResult> OnPostAsync()
        {
            if (!ModelState.IsValid)
            {
                return Page();
            }

            _context.Todos.Add(Todo);
            await _context.SaveChangesAsync();

            return RedirectToPage("/Index");
        }
    }
}