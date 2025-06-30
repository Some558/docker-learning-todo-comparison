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
            return Page();
        }

        [BindProperty]
        public Todo Todo { get; set; } = default!;

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